import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { PageHeader, Button, Card, Modal } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface Customer {
  id: string;
  country_code: string;
  company_name: string;
  trading_name: string;
  account_number: string;
  email: string;
  credit_limit: number;
  credit_balance: number;
  is_active: boolean;
  is_credit_hold: boolean;
  created_at: string;
}

interface PagedResult {
  items: Customer[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export default function CustomersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [showCreate, setShowCreate] = useState(false);

  const { data, isLoading } = useQuery<PagedResult>({
    queryKey: ['customers', page, search],
    queryFn: () => apiClient.get(`/customer/customers?page=${page}&page_size=20&search=${search}`).then((r) => r.data),
  });

  return (
    <div>
      <PageHeader
        title="Customers"
        subtitle="Manage customer accounts"
        actions={<Button onClick={() => setShowCreate(true)}>Create Customer</Button>}
      />

      <div className="mt-4 mb-4">
        <input
          type="text"
          placeholder="Search customers..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="w-full max-w-md px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <Card className="mt-2">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Company</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Account #</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Country</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Credit</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {isLoading ? (
                <tr><td colSpan={6} className="px-6 py-4 text-center text-gray-500">Loading...</td></tr>
              ) : data?.items?.length ? (
                data.items.map((c) => (
                  <tr key={c.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm">
                      <Link to={`/customers/${c.id}`} className="font-medium text-blue-600 hover:text-blue-800">{c.company_name}</Link>
                      {c.trading_name && <span className="block text-xs text-gray-400">t/a {c.trading_name}</span>}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500 font-mono">{c.account_number}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{c.country_code}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{c.email}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">${c.credit_balance.toFixed(2)} / ${c.credit_limit.toFixed(2)}</td>
                    <td className="px-6 py-4">
                      <div className="flex gap-1">
                        <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${c.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                          {c.is_active ? 'Active' : 'Inactive'}
                        </span>
                        {c.is_credit_hold && (
                          <span className="inline-flex px-2 py-1 text-xs font-medium rounded-full bg-yellow-100 text-yellow-700">Hold</span>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={6} className="px-6 py-4 text-center text-gray-500">No customers found</td></tr>
              )}
            </tbody>
          </table>
        </div>

        {data && data.total_pages > 1 && (
          <div className="flex items-center justify-between px-6 py-3 border-t border-gray-200">
            <span className="text-sm text-gray-700">Page {data.page} of {data.total_pages} ({data.total_count} total)</span>
            <div className="flex space-x-2">
              <Button variant="secondary" size="sm" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1}>Previous</Button>
              <Button variant="secondary" size="sm" onClick={() => setPage(p => p + 1)} disabled={page >= data.total_pages}>Next</Button>
            </div>
          </div>
        )}
      </Card>

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create Customer">
        <p className="text-sm text-gray-500">Customer creation form coming soon.</p>
      </Modal>
    </div>
  );
}
