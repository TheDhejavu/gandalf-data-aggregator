import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useToastMessage } from '../context/error-context';

const Loader: React.FC = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      style={{ margin: 'auto',marginTop: -20, background: '#f7f7f7', display: 'block' }}
      width="100px"
      height="100px"
      viewBox="0 0 100 100"
      preserveAspectRatio="xMidYMid"
    >
      <defs>
        <clipPath id="ldio-oa510johls-cp" x="0" y="0" width="100" height="100">
          <path d="M81.3,58.7H18.7c-4.8,0-8.7-3.9-8.7-8.7v0c0-4.8,3.9-8.7,8.7-8.7h62.7c4.8,0,8.7,3.9,8.7,8.7v0C90,54.8,86.1,58.7,81.3,58.7z"></path>
        </clipPath>
      </defs>
      <path fill="none" stroke="#575757" strokeWidth="2.7928" d="M82 63H18c-7.2,0-13-5.8-13-13v0c0-7.2,5.8-13,13-13h64c7.2,0,13,5.8,13,13v0C95,57.2,89.2,63,82,63z"></path>
      <g clipPath="url(#ldio-oa510johls-cp)">
        <g>
          <rect x="-100" y="0" width="25" height="100" fill="#e15b64"></rect>
          <rect x="-75" y="0" width="25" height="100" fill="#f47e60"></rect>
          <rect x="-50" y="0" width="25" height="100" fill="#f8b26a"></rect>
          <rect x="-25" y="0" width="25" height="100" fill="#abbd81"></rect>
          <rect x="0" y="0" width="25" height="100" fill="#e15b64"></rect>
          <rect x="25" y="0" width="25" height="100" fill="#f47e60"></rect>
          <rect x="50" y="0" width="25" height="100" fill="#f8b26a"></rect>
          <rect x="75" y="0" width="25" height="100" fill="#abbd81"></rect>
          <animateTransform attributeName="transform" type="translate" dur="1s" repeatCount="indefinite" keyTimes="0;1" values="0;100"></animateTransform>
        </g>
      </g>
    </svg>
  );
};

const Auth: React.FC = () => {
  const [responseData, _] = useState<any>(null);
  const authURL = window.location.href; 
  const { showToast } = useToastMessage();
  const navigate = useNavigate();

  const completeAuth = async () => {
    try {
      const response = await fetch(`${import.meta.env.VITE_SERVER_URL}/auth/twitter/complete`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          authURL: authURL 
        })
      });

      if (!response.ok) {
        const errorMessage = await response.json();
        throw new Error(errorMessage["message"]);
      }

      const data = await response.json();
      localStorage.setItem('token', data.token);

      navigate('/'); 
      showToast(`Welcome ${data["username"]}, let's get you started!`, { type: 'success' }); 
    } catch (error: any) {
        navigate('/login'); 
        showToast(error.message, { type: 'error' });
    }
  };

  useEffect(() => {
    completeAuth(); 
  }, []);

  return (
    <div className='flex item-center'>
      {responseData ? (
        <div>
          <h2>Response Data:</h2>
          <pre>{JSON.stringify(responseData, null, 2)}</pre>
        </div>
      ) : (<> 
       <div className="h-screen flex flex-col items-center justify-center w-full">
        <div className="mx-auto text-center">
            <p className="text-sm text-gray-700">
                Wait while we connect your twitter accont...
            </p>
            <Loader />
        </div>
        </div>
      </>)}
    </div>
  );
};

export default Auth;
