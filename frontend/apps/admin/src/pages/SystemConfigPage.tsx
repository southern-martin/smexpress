import { useQuery } from '@tanstack/react-query';
import { PageHeader, Button, Card, Modal } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';
import { useState } from 'react';

interface SystemConfig {
  id: string;
  key: string;
  value: string;
  description: string;
  country_code: string;
}

export default function SystemConfigPage() {
  const [showCreate, setShowCreate] = useState(false);

  const { data, isLoading } = useQuery<SystemConfig[]>({
    queryKey: ['system-configs'],
    queryFn: () => apiClient.get('/config/system-configs').then((r) => r.data),
  });

  return (
    <div>
      <PageHeader
        title="System Configuration"
        subtitle="Global and country-specific settings"
        actions={<Button onClick={() => setShowCreate(true)}>Add Config</Button>}
      />

      <Card className="mt-6">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Key</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Value</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Country</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {isLoading ? (
                <tr><td colSpan={4} className="px-6 py-4 text-center text-gray-500">Loading...</td></tr>
              ) : data?.length ? (
                data.map((cfg) => (
                  <tr key={cfg.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-medium text-gray-900 font-mono">{cfg.key}</td>
                    <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">{cfg.value}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{cfg.country_code || 'Global'}</td>
                    <td className="px-6 py-4 text-sm text-gray-500">{cfg.description}</td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={4} className="px-6 py-4 text-center text-gray-500">No configs found</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Add Configuration">
        <p className="text-sm text-gray-500">Config form coming soon.</p>
      </Modal>
    </div>
  );
}
