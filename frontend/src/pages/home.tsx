import React, { useState , useEffect } from 'react';
import Pagination from '../components/pagination';
import Modal from '../components/modal';
import { Line } from 'react-chartjs-2';
import { toast } from 'react-toastify';
import Loader from '../components/loader';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
} from 'chart.js';
import Placeholder from '../components/placeholder';

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);


function formatDate(dateString: string) {
    const date = new Date(dateString);
    const formattedDate = date.toLocaleDateString('en-GB');
    return formattedDate;
}

const Home: React.FC= () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isActivityLoading, setActivityLoading] = useState(true);
  const toggleModal = () => setIsModalOpen(!isModalOpen);

  const [selectedYear, setSelectedYear] = useState<string>(new Date().getFullYear().toString());
  const [stats, setStats] = useState<{ year_data: Record<string, { Labels: number[], Months: string[] }>, current_year: string } | null>(null);
  const yearlyData = stats && stats.year_data[selectedYear] ? stats.year_data[selectedYear].Labels : [];
  const updateChartData = (year: string) => {
    setSelectedYear(year);
  };

  const labels = stats && stats.year_data[selectedYear] ? stats.year_data[selectedYear].Months : [];

  const data = {
    labels,
    datasets: [{
      label: 'Movies/Series Watched',
      data: yearlyData || [], 
      backgroundColor: '#FF6E33',
      borderColor: '#DF3F10',
      borderWidth: 1
    }]
  };

  const options = {
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  const [responseActivityData, setResponseActivityData] = useState<any>({});
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [startIndex, setStartIndex] = useState(0);
  
  const limit = 10

  const fetchActivityData = async (page: number, limit: number) => {
    try {
     
      const storedToken = localStorage.getItem('token');
      console.log(storedToken)
      const response = await fetch(`${process.env.REACT_APP_SERVER_URL}/user/activity?limit=${limit}&page=${page}`, {
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

      const jsonResponse = await response.json();
      setResponseActivityData(jsonResponse.activities);
      setTotalPages(jsonResponse.activities.total);

      // set stats
      setStats(jsonResponse.stats);
      setSelectedYear(jsonResponse.stats.current_year);
  
      setStartIndex((page - 1) * limit + 1);
      setActivityLoading(false)

    } catch (error: any) {
      // showToast(error.message, { type: 'error' });
      toast(error.message, { type: 'error' });
      setActivityLoading(false)
    }
  };

  useEffect(() => {
    fetchActivityData(currentPage, limit); 
  }, [currentPage]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const [user, setUser] = useState<any>({
    avatar_url: "https://as1.ftcdn.net/v2/jpg/05/16/27/58/1000_F_516275801_f3Fsp17x6HQK0xQgDQEELoTuERO4SsWV.jpg",
    username: "loading....."
  });
  useEffect(() => {
    const currentUser = async () => {
        try {
            const storedToken = localStorage.getItem('token');
            const response = await fetch(`${process.env.REACT_APP_SERVER_URL}/user/me`, {
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
           
            setUser(data)
        } catch (error: any) {
          toast(error.message, { type: 'error' });
        }
    };
    
    currentUser(); 
}, []); 

return (
    <>
    <Modal isOpen={isModalOpen} toggleModal={toggleModal} />
    <div className="bg-card w-full max-w-3xl mx-auto">
      <div className="">
        <div className="py-2 px-2 border-b">
            <div className="flex font-semibold whitespace-nowrap leading-none tracking-tight">
                <div className="pt-2 pr-2 p-1">
                <svg width="34px" height="34px" strokeWidth="1.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" color="#000000">
                    <path d="M12 11.5C14.2091 11.5 16 9.70914 16 7.5C16 5.29086 14.2091 3.5 12 3.5C9.79086 3.5 8 5.29086 8 7.5C8 9.70914 9.79086 11.5 12 11.5Z" stroke="#000000" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"></path>
                    <path d="M7 20.5C9.20914 20.5 11 18.7091 11 16.5C11 14.2909 9.20914 12.5 7 12.5C4.79086 12.5 3 14.2909 3 16.5C3 18.7091 4.79086 20.5 7 20.5Z" stroke="#000000" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"></path>
                    <path d="M17 20.5C19.2091 20.5 21 18.7091 21 16.5C21 14.2909 19.2091 12.5 17 12.5C14.7909 12.5 13 14.2909 13 16.5C13 18.7091 14.7909 20.5 17 20.5Z" stroke="#000000" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"></path>
                </svg>
                </div>
                <h2 className="pt-3 text-lg sm:text-xl">Data Aggregator</h2>
            </div>
            <p className="text-xs pl-2 text-muted-foreground text-gray-500 dark:text-gray-400">
                Welcome!, Let's see what you have been up to.
            </p>
        </div>
        <div className="relative m-3 flex flex-col space-y-1.5 ">
            <div className="flex max-w-40 p-2">
                <div className="pr-2">
                    <img src={user && user.avatar_url} alt="User avatar" className="rounded-full  object-cover" width={35} height={35} />
                </div>
                <div>
                    <h3 className="text-lg text-gray-800 text-sm">{user && user.username}</h3>
                    <p className="text-xs text-muted-foreground text-gray-500 dark:text-gray-400">
                        #100.
                    </p>
                </div>
            </div>
            <div className="absolute right-0 mr-4 flex justify-end">
                <button onClick={toggleModal}  style={{ borderRadius: '100px', }} className="inline-flex items-center justify-center whitespace-nowrap rounded-2xl text-sm font-medium ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-black text-white hover:bg-black/70 transition-all px-4 py-2">
                    Connect Netflix
                </button>
            </div>
        </div>  
    </div>
      <div className="p-0 mt-4">
        <div className="grid w-full overflow-auto">
          <div className="relative w-full overflow-auto">
          {(responseActivityData.data && responseActivityData.data.length > 0 ?
           <>
            <table className="w-full caption-bottom text-sm">
                <thead className="[&amp;_tr]:border-b">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                        <th className="h-12 px-4 text-left align-middle text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 font-bold">S/N</th>
                        <th className="h-12 px-4 text-left align-middle  text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 font-bold">Movie Title</th>
                        <th className="h-12 px-4 text-left align-middle  text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 font-bold">Activity Date</th>
                        <th className="h-12 px-4 text-left align-middle  text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 font-bold">IMDB</th>
                    </tr>
              </thead>
              <tbody>
                {responseActivityData.data &&
                    responseActivityData.data.map((item: any, index: number) => (
                    <tr key={index}>
                        <td className="p-4 align-middle font-semibold">{startIndex + index}</td>
                        <td className="p-4 align-middle font-medium">{item.title}</td>
                        <td className="p-4 align-middle">{formatDate(item.date)}</td>
                        <td className="p-4 align-middle">
                          <a className="text-sm text-gray-800 underline" href={`https://imdb.com/title/${item.subject.find((subject: any) => subject.identifier_type === "IMDB")?.value}`} target='_black'>
                            {(item.subject.find((subject: any) => subject.identifier_type === "IMDB")?.value)}
                          </a>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
           
            {responseActivityData.data && responseActivityData.data.length > 0 && (
                <Pagination total={totalPages} limit={limit}  currentPage={currentPage} onPageChange={handlePageChange} />
            )}
            </>
            : 
            <Placeholder.Placeholder loading={isActivityLoading} text="Sorry, there is no data to display." toggleModal={toggleModal}></Placeholder.Placeholder>
        )}
          </div>
        </div>
    
        <div className="p-10 border rounded-lg bg-white" hidden={isActivityLoading}>
          <h1 className="">Netflix Usage Analytics</h1>
          {(stats && Object.keys(stats.year_data).length > 0 ?
          <>
          <div>
            <label htmlFor="yearSelect">Select Year:</label>
            <select id="yearSelect" onChange={(e) => updateChartData(e.target.value)} value={selectedYear}>
            {stats && Object.keys(stats.year_data).map((year) => (
              <option key={year} value={year}>{year}</option>
            ))}
            </select>
          </div>
          <Line data={data} options={options} />
          </>
          :  <Placeholder.AnalyticPlaceholder ></Placeholder.AnalyticPlaceholder> 
        )}
        </div>
      </div>
    </div>
    <footer className="bottom-0 w-full py-4 text-center text-xs text-gray-600 dark:text-gray-400">
        © 2024 Data Aggregator. All rights reserved.
    </footer></>
  );
};

export default Home;
