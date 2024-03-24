
const Login: React.FC = () => {

  return (
    <>
    <div className="text-black h-screen flex flex-col justify-center items-center py-12 space-y-6 text-center md:space-y-5">
        <div className='bg-white p-20 border rounded-lg '>
            <div className="space-y-2 bg-0">
                <div className="flex items-center justify-center space-x-4">
                    <svg width="34px" height="34px" strokeWidth="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#000000"><path d="M12 11.5C14.2091 11.5 16 9.70914 16 7.5C16 5.29086 14.2091 3.5 12 3.5C9.79086 3.5 8 5.29086 8 7.5C8 9.70914 9.79086 11.5 12 11.5Z" stroke="#000000" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"></path><path d="M7 20.5C9.20914 20.5 11 18.7091 11 16.5C11 14.2909 9.20914 12.5 7 12.5C4.79086 12.5 3 14.2909 3 16.5C3 18.7091 4.79086 20.5 7 20.5Z" stroke="#000000" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"></path><path d="M17 20.5C19.2091 20.5 21 18.7091 21 16.5C21 14.2909 19.2091 12.5 17 12.5C14.7909 12.5 13 14.2909 13 16.5C13 18.7091 14.7909 20.5 17 20.5Z" stroke="#000000" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"></path></svg>
                    <h1 className="text-4xl font-bold">Data Aggregator</h1>
                </div>
                <p className="text-sm text-gray-700 pt-2">
                Data aggregator allows you to get insights into your digital footprints.
                </p>
            </div>
            <div className="w-full my-4 space-y-2">
                <div>
                <a   style={{ borderRadius: '100px', }} href={`${process.env.REACT_APP_SERVER_URL}/auth/twitter`}  target="_blank" className="bg-black text-white block hover:text-gray-400 text-black  hover:bg-gray-900 transition-colors  text-sm font-medium px-2 py-3 w-full focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50">
                    Sign in with X
                </a>
                </div>
            </div>
        </div>
    </div>
    </>
  );
};

export default Login;
