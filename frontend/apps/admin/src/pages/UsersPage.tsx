import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { PageHeader, Button, Card, Modal } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  country_code: string;
  is_active: boolean;
  created_at: string;
}

interface PagedResult {
  items: User[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export default function UsersPage() {
  const [page, setPage] = useState(1);
  const [showCreate, setShowCreate] = useState(false);

  const { data, isLoading } = useQuery<PagedResult>({
    queryKey: ['users', page],
    queryFn: () => apiClient.get(`/auth/users?page=${page}&page_size=20`).then((r) => r.data),
  });

  return (
    <div>
      <PageHeader
        title="Users"
        subtitle="Manage system users"
        actions={<Button onClick={() => setShowCreate(true)}>Create User</Button>}
      />

      <Card className="mt-6">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Country</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {isLoading ? (
                <tr><td colSpan={5} className="px-6 py-4 text-center text-gray-500">Loading...</td></tr>
              ) : data?.items?.length ? (
                data.items.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">{user.first_name} {user.last_name}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{user.email}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{user.country_code}</td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${user.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                        {user.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">{new Date(user.created_at).toLocaleDateString()}</td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={5} className="px-6 py-4 text-center text-gray-500">No users found</td></tr>
              )}
            </tbody>
          </table>
        </div>

        {data && data.total_pages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 border-t border-gray-200">
            <span className="text-sm text-gray-700">
              Page {data.page} of {data.total_pages} ({data.total_count} total)
            </span>
            <div className="flex space-x-2">
              <Button variant="secondary" size="sm" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1}>Previous</Button>
              <Button variant="secondary" size="sm" onClick={() => setPage(p => p + 1)} disabled={page >= data.total_pages}>Next</Button>
            </div>
          </div>
        )}
      </Card>

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create User">
        <p className="text-sm text-gray-500">User creation form coming soon.</p>
      </Modal>
    </div>
  );
}
