import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { PageHeader, Card, Button } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';

interface PostcodeResult {
  id: string;
  postcode: string;
  suburb: string;
  city: string;
  state: string;
  state_code: string;
}

export default function AddressLookupPage() {
  const [countryCode, setCountryCode] = useState('AU');
  const [query, setQuery] = useState('');
  const [searchTerm, setSearchTerm] = useState('');

  const { data, isLoading } = useQuery<{ items: PostcodeResult[]; total_count: number }>({
    queryKey: ['postcode-search', countryCode, searchTerm],
    queryFn: () => apiClient.get(`/address/postcodes/search?country_code=${countryCode}&q=${searchTerm}&page_size=50`).then((r) => r.data),
    enabled: searchTerm.length >= 2,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchTerm(query);
  };

  return (
    <div>
      <PageHeader title="Address Lookup" subtitle="Search postcodes and suburbs" />

      <Card className="mt-6" padding>
        <form onSubmit={handleSearch} className="flex gap-3 items-end">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Country</label>
            <select
              value={countryCode}
              onChange={(e) => setCountryCode(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="AU">Australia</option>
              <option value="NZ">New Zealand</option>
              <option value="UK">United Kingdom</option>
              <option value="CA">Canada</option>
              <option value="US">United States</option>
            </select>
          </div>
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-1">Postcode or Suburb</label>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="e.g. 2000 or Sydney"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <Button type="submit">Search</Button>
        </form>
      </Card>

      {searchTerm && (
        <Card className="mt-4">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Postcode</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Suburb</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">City</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">State</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {isLoading ? (
                  <tr><td colSpan={4} className="px-6 py-4 text-center text-gray-500">Searching...</td></tr>
                ) : data?.items?.length ? (
                  data.items.map((p) => (
                    <tr key={p.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 text-sm font-mono text-gray-900">{p.postcode}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">{p.suburb}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">{p.city}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">{p.state} {p.state_code && `(${p.state_code})`}</td>
                    </tr>
                  ))
                ) : (
                  <tr><td colSpan={4} className="px-6 py-4 text-center text-gray-500">No results found</td></tr>
                )}
              </tbody>
            </table>
          </div>
          {data && <p className="px-6 py-2 text-xs text-gray-400">{data.total_count} results</p>}
        </Card>
      )}
    </div>
  );
}
