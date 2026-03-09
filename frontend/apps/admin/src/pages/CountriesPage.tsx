import { useQuery } from '@tanstack/react-query';
import { PageHeader, Card } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface CountryConfig {
  id: string;
  country_code: string;
  country_name: string;
  currency_code: string;
  timezone: string;
  is_active: boolean;
}

export default function CountriesPage() {
  const { data, isLoading } = useQuery<CountryConfig[]>({
    queryKey: ['countries'],
    queryFn: () => apiClient.get('/config/countries').then((r) => r.data),
  });

  return (
    <div>
      <PageHeader title="Countries" subtitle="Country configurations" />

      <Card className="mt-6">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Code</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Currency</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Timezone</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {isLoading ? (
                <tr><td colSpan={5} className="px-6 py-4 text-center text-gray-500">Loading...</td></tr>
              ) : data?.length ? (
                data.map((c) => (
                  <tr key={c.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">{c.country_code}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{c.country_name}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{c.currency_code}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{c.timezone}</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${c.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                        {c.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={5} className="px-6 py-4 text-center text-gray-500">No countries configured</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
