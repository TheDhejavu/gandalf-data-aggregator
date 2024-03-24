import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import Loader from '../components/loader';

const Auth: React.FC = () => {
  const [responseData, setResponseData] = useState<any>(null);
  const authURL = window.location.href; 
 
  const navigate = useNavigate();

  const completeAuth = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_SERVER_URL}/auth/twitter/complete`, {
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

      window.location.href = '/';
      toast(`Welcome ${data["username"]}, let's get you started!`, { type: 'success' }); 
    } catch (error: any) {
      navigate('/login'); 
      toast(error.message, { type: 'error' });
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
       <div className="h-screen flex items-center justify-center w-full">
        <div className="mx-auto text-center">
            <p className="text-sm">
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
