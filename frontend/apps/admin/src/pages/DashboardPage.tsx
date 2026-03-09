import { PageHeader } from '@smexpress/ui';

export default function DashboardPage() {
  const stats = [
    { label: 'Total Users', value: '-' },
    { label: 'Active Franchises', value: '-' },
    { label: 'Countries', value: '-' },
    { label: 'Pending Withdrawals', value: '-' },
  ];

  return (
    <div>
      <PageHeader title="Dashboard" subtitle="Overview of your system" />
      <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <div key={stat.label} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <p className="text-sm text-gray-500">{stat.label}</p>
            <p className="text-3xl font-bold text-gray-900 mt-1">{stat.value}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
