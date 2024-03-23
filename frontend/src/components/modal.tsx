import Connect  from '@gandalf-network/connect';
import { useToastMessage } from '../context/error-context';
import { useEffect, useState } from 'react';

interface ModalProps {
  isOpen: boolean;
  toggleModal: () => void;
}

const Modal: React.FC<ModalProps> = ({ isOpen, toggleModal }) => {
    const [redirectURL, setRedirectURL] = useState('');
    const [qrCodeDataUrl, setQrCodeDataUrl] = useState<string>('');
    const [connectURL, setConnectURL] = useState<string>('');
    const { showToast } = useToastMessage();

    useEffect(() => {
        const generateRedirectURL = async () => {
            try {
                const storedToken = localStorage.getItem('token');
                const response = await fetch(`${process.env.REACT_APP_SERVER_URL}/user/generate-callback`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${storedToken}` 
                    },
                });
        
                if (!response.ok) {
                    const errorMessage = await response.json();
                    throw new Error(errorMessage["message"]);
                }
        
                const data = await response.json();
                setRedirectURL(data.callbackURL);
            } catch (error: any) {
                showToast(error.message, { type: 'error' });
            }
        };
        
        generateRedirectURL(); 
    }, []); 

    useEffect(() => {
        const connect = new Connect({
            publicKey: '0x02290a12d7ed02b3377e683d3285d6e4282662fc6f6d7496754d4d01e3085e6452',
            redirectURL: redirectURL, 
            services: { 'netflix': true }
        });

        const fetchConnectData = async () => {
            try {
                const url = await connect.generateURL();
                const qrCode = await connect.generateQRCode();
                console.log(url)
                setConnectURL(url);
                setQrCodeDataUrl(qrCode);

            } catch (error) {
                console.error('Error generating connect data:', error);
            }
        };

        if (redirectURL) {
          fetchConnectData();
        }
    }, [redirectURL]); 

  return (
    <>
      {isOpen && (
      <div className="modal opacity-1 fixed w-full h-full top-0 left-0 flex items-center justify-center z-10">
        <div className="modal-overlay absolute w-full h-full bg-gray-900 opacity-50 "></div>

        <div onClick={toggleModal} className="modal-close absolute z-0 top-0 right-0 cursor-pointer flex flex-col items-center mt-4 mr-4 text-white text-sm z-0">
            <svg className="fill-current text-white" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 18 18">
              <path d="M14.53 4.53l-1.06-1.06L9 7.94 4.53 3.47 3.47 4.53 7.94 9l-4.47 4.47 1.06 1.06L9 10.06l4.47 4.47 1.06-1.06L10.06 9z"></path>
            </svg>
            <span className="text-sm">(Esc)</span>
        </div>

        <div className="modal-container bg-white w-11/12 md:max-w-md mx-auto rounded shadow-xl p-6 mb-6 relative">
          <div className="modal-content text-center relative p-6">
            <div className="">
              <p className="text-2xl font-bold text-center">Connect to Netflix</p>
              <div className="absolute z-0 right-0 top-0 mt-4  mr-4 modal-close cursor-pointer" onClick={toggleModal}>
                <svg className="fill-current text-black" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 18 18">
                  <path d="M14.53 4.53l-1.06-1.06L9 7.94 4.53 3.47 3.47 4.53 7.94 9l-4.47 4.47 1.06 1.06L9 10.06l4.47 4.47 1.06-1.06L10.06 9z"></path>
                </svg>
              </div>
            </div>
            <p className="text-sm py-2 text-center text-muted-foreground text-gray-500 dark:text-gray-400">
              Scan the QRCode below to connect your netflix account
            </p>
            <div className="flex justify-center items-center">
              <img src={qrCodeDataUrl} alt="Connect QRCode" />
            </div>
            <a href={connectURL} target="_blank" className="hover:text-gray-400 transition-colors  text-sm py-2 text-center text-muted-foreground text-gray-500"> Using a mobile device? Tap here! </a>
          </div>
        </div>
      </div>
     )}
     </>
  );
};


export default Modal;