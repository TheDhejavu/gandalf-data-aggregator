import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { GlobalProvider } from './context/error-context';
import {  ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <GlobalProvider>
    <ToastContainer />
    <App />
  </GlobalProvider>
);
