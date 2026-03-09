import { useQuery } from '@tanstack/react-query';
import { PageHeader, Button, Card, Modal } from '@smexpress/ui';
import apiClient from '@smexpress/api-client';
import { useState } from 'react';

interface Role {
  id: string;
  name: string;
  display_name: string;
  description: string;
  is_system: boolean;
  permissions: string[];
}

export default function RolesPage() {
  const [showCreate, setShowCreate] = useState(false);

  const { data, isLoading } = useQuery<Role[]>({
    queryKey: ['roles'],
    queryFn: () => apiClient.get('/auth/roles').then((r) => r.data),
  });

  return (
    <div>
      <PageHeader
        title="Roles"
        subtitle="Manage roles and permissions"
        actions={<Button onClick={() => setShowCreate(true)}>Create Role</Button>}
      />

      <div className="mt-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {isLoading ? (
          <p className="text-gray-500 col-span-full text-center py-8">Loading...</p>
        ) : data?.length ? (
          data.map((role) => (
            <Card key={role.id} padding>
              <div className="flex items-center justify-between mb-2">
                <h3 className="text-lg font-semibold text-gray-900">{role.display_name}</h3>
                {role.is_system && (
                  <span className="px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded">System</span>
                )}
              </div>
              <p className="text-sm text-gray-500 mb-3">{role.description}</p>
              {role.permissions?.length > 0 && (
                <div className="flex flex-wrap gap-1">
                  {role.permissions.slice(0, 5).map((p) => (
                    <span key={p} className="px-2 py-0.5 text-xs bg-blue-50 text-blue-700 rounded">{p}</span>
                  ))}
                  {role.permissions.length > 5 && (
                    <span className="px-2 py-0.5 text-xs bg-gray-100 text-gray-500 rounded">+{role.permissions.length - 5} more</span>
                  )}
                </div>
              )}
            </Card>
          ))
        ) : (
          <p className="text-gray-500 col-span-full text-center py-8">No roles found</p>
        )}
      </div>

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create Role">
        <p className="text-sm text-gray-500">Role creation form coming soon.</p>
      </Modal>
    </div>
  );
}
