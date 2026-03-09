import { Routes, Route } from 'react-router-dom';

function Home() {
  return <div className="p-8"><h1 className="text-2xl font-bold">SmExpress Customer Portal</h1></div>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
    </Routes>
  );
}
