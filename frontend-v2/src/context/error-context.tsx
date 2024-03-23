import React, { createContext, useContext, useState, ReactNode } from 'react';
import { toast, ToastOptions } from 'react-toastify';

interface GlobalContextType {
  message: string | null;
  showToast: (message: string, options?: ToastOptions) => void;
}

const GlobalContext = createContext<GlobalContextType>({
  message: null,
  showToast: () => {},
});

export const useToastMessage = () => useContext(GlobalContext);

export const GlobalProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [message, _] = useState<string | null>(null);

  const showToast = (message: string, options?: ToastOptions) => {
    toast(message, options);
  };

  return (
    <GlobalContext.Provider value={{ message, showToast }}>
      {children}
    </GlobalContext.Provider>
  );
};