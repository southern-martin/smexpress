import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { PageHeader, Button, Card, Modal } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface Franchise {
  id: string;
  country_code: string;
  name: string;
  code: string;
  contact_name: string;
  email: string;
  is_active: boolean;
  commission_rate: number;
  created_at: string;
}

interface PagedResult {
  items: Franchise[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export default function FranchisesPage() {
  const [page, setPage] = useState(1);
  const [showCreate, setShowCreate] = useState(false);

  const { data, isLoading } = useQuery<PagedResult>({
    queryKey: ['franchises', page],
    queryFn: () => apiClient.get(`/franchise/franchises?page=${page}&page_size=20`).then((r) => r.data),
  });

  return (
    <div>
      <PageHeader
        title="Franchises"
        subtitle="Manage franchise partners"
        actions={<Button onClick={() => setShowCreate(true)}>Create Franchise</Button>}
      />

      <Card className="mt-6">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Code</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Country</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Contact</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Commission</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {isLoading ? (
                <tr><td colSpan={6} className="px-6 py-4 text-center text-gray-500">Loading...</td></tr>
              ) : data?.items?.length ? (
                data.items.map((f) => (
                  <tr key={f.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">{f.name}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{f.code}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{f.country_code}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{f.contact_name}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{(f.commission_rate * 100).toFixed(1)}%</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${f.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                        {f.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={6} className="px-6 py-4 text-center text-gray-500">No franchises found</td></tr>
              )}
            </tbody>
          </table>
        </div>

        {data && data.total_pages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 border-t border-gray-200">
            <span className="text-sm text-gray-700">Page {data.page} of {data.total_pages}</span>
            <div className="flex space-x-2">
              <Button variant="secondary" size="sm" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1}>Previous</Button>
              <Button variant="secondary" size="sm" onClick={() => setPage(p => p + 1)} disabled={page >= data.total_pages}>Next</Button>
            </div>
          </div>
        )}
      </Card>

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create Franchise">
        <p className="text-sm text-gray-500">Franchise creation form coming soon.</p>
      </Modal>
    </div>
  );
}
