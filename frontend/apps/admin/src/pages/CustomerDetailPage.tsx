import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { PageHeader, Card, Button } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface Customer {
  id: string;
  country_code: string;
  company_name: string;
  trading_name: string;
  account_number: string;
  abn: string;
  email: string;
  phone: string;
  website: string;
  credit_limit: number;
  credit_balance: number;
  payment_terms: number;
  is_active: boolean;
  is_credit_hold: boolean;
}

interface Contact {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  position: string;
  is_primary: boolean;
}

interface Address {
  id: string;
  address_type: string;
  company_name: string;
  contact_name: string;
  address_line1: string;
  address_line2: string;
  city: string;
  state: string;
  postcode: string;
  country_code: string;
  is_default: boolean;
}

export default function CustomerDetailPage() {
  const { id } = useParams<{ id: string }>();

  const { data: customer, isLoading } = useQuery<Customer>({
    queryKey: ['customer', id],
    queryFn: () => apiClient.get(`/customer/customers/${id}`).then((r) => r.data),
  });

  const { data: contacts } = useQuery<Contact[]>({
    queryKey: ['customer-contacts', id],
    queryFn: () => apiClient.get(`/customer/customers/${id}/contacts`).then((r) => r.data),
    enabled: !!id,
  });

  const { data: addresses } = useQuery<Address[]>({
    queryKey: ['customer-addresses', id],
    queryFn: () => apiClient.get(`/customer/customers/${id}/addresses`).then((r) => r.data),
    enabled: !!id,
  });

  if (isLoading) return <div className="p-8 text-gray-500">Loading...</div>;
  if (!customer) return <div className="p-8 text-gray-500">Customer not found</div>;

  return (
    <div>
      <PageHeader
        title={customer.company_name}
        subtitle={`Account: ${customer.account_number || 'N/A'} | ${customer.country_code}`}
        actions={
          <Link to="/customers">
            <Button variant="secondary">Back to List</Button>
          </Link>
        }
      />

      <div className="mt-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Details */}
        <Card padding>
          <h3 className="text-lg font-semibold mb-4">Details</h3>
          <dl className="space-y-2 text-sm">
            <div><dt className="text-gray-500">Trading Name</dt><dd>{customer.trading_name || '-'}</dd></div>
            <div><dt className="text-gray-500">ABN</dt><dd>{customer.abn || '-'}</dd></div>
            <div><dt className="text-gray-500">Email</dt><dd>{customer.email || '-'}</dd></div>
            <div><dt className="text-gray-500">Phone</dt><dd>{customer.phone || '-'}</dd></div>
            <div><dt className="text-gray-500">Website</dt><dd>{customer.website || '-'}</dd></div>
            <div><dt className="text-gray-500">Payment Terms</dt><dd>{customer.payment_terms} days</dd></div>
            <div><dt className="text-gray-500">Credit Limit</dt><dd>${customer.credit_limit.toFixed(2)}</dd></div>
            <div><dt className="text-gray-500">Credit Balance</dt><dd>${customer.credit_balance.toFixed(2)}</dd></div>
            <div className="flex gap-2 pt-2">
              <span className={`px-2 py-1 text-xs rounded-full ${customer.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                {customer.is_active ? 'Active' : 'Inactive'}
              </span>
              {customer.is_credit_hold && (
                <span className="px-2 py-1 text-xs rounded-full bg-yellow-100 text-yellow-700">Credit Hold</span>
              )}
            </div>
          </dl>
        </Card>

        {/* Contacts */}
        <Card padding>
          <h3 className="text-lg font-semibold mb-4">Contacts</h3>
          {contacts?.length ? (
            <ul className="space-y-3">
              {contacts.map((c) => (
                <li key={c.id} className="border-b border-gray-100 pb-2 last:border-0">
                  <p className="font-medium text-sm">
                    {c.first_name} {c.last_name}
                    {c.is_primary && <span className="ml-2 text-xs bg-blue-100 text-blue-700 px-1 rounded">Primary</span>}
                  </p>
                  {c.position && <p className="text-xs text-gray-400">{c.position}</p>}
                  {c.email && <p className="text-xs text-gray-500">{c.email}</p>}
                  {c.phone && <p className="text-xs text-gray-500">{c.phone}</p>}
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-gray-400">No contacts</p>
          )}
        </Card>

        {/* Addresses */}
        <Card padding>
          <h3 className="text-lg font-semibold mb-4">Addresses</h3>
          {addresses?.length ? (
            <ul className="space-y-3">
              {addresses.map((a) => (
                <li key={a.id} className="border-b border-gray-100 pb-2 last:border-0">
                  <p className="font-medium text-sm">
                    {a.company_name || a.contact_name}
                    {a.is_default && <span className="ml-2 text-xs bg-blue-100 text-blue-700 px-1 rounded">Default</span>}
                    <span className="ml-2 text-xs bg-gray-100 text-gray-500 px-1 rounded">{a.address_type}</span>
                  </p>
                  <p className="text-xs text-gray-500">{a.address_line1}</p>
                  <p className="text-xs text-gray-500">{a.city} {a.state} {a.postcode} {a.country_code}</p>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-gray-400">No addresses</p>
          )}
        </Card>
      </div>
    </div>
  );
}
