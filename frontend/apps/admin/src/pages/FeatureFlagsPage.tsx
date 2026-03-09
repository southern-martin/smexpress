import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { PageHeader, Card } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface FeatureFlag {
  id: string;
  name: string;
  key: string;
  description: string;
  is_enabled: boolean;
  country_code: string;
}

export default function FeatureFlagsPage() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery<FeatureFlag[]>({
    queryKey: ['feature-flags'],
    queryFn: () => apiClient.get('/config/feature-flags').then((r) => r.data),
  });

  const toggleMutation = useMutation({
    mutationFn: (id: string) => apiClient.patch(`/config/feature-flags/${id}/toggle`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['feature-flags'] }),
  });

  return (
    <div>
      <PageHeader title="Feature Flags" subtitle="Toggle features across countries" />

      <Card className="mt-6">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Key</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Country</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Enabled</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {isLoading ? (
                <tr><td colSpan={5} className="px-6 py-4 text-center text-gray-500">Loading...</td></tr>
              ) : data?.length ? (
                data.map((flag) => (
                  <tr key={flag.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">{flag.name}</td>
                    <td className="px-6 py-4 text-sm text-gray-500 font-mono">{flag.key}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{flag.country_code || 'All'}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{flag.description}</td>
                    <td className="px-6 py-4">
                      <button
                        onClick={() => toggleMutation.mutate(flag.id)}
                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${flag.is_enabled ? 'bg-blue-600' : 'bg-gray-300'}`}
                      >
                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${flag.is_enabled ? 'translate-x-6' : 'translate-x-1'}`} />
                      </button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={5} className="px-6 py-4 text-center text-gray-500">No feature flags</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
