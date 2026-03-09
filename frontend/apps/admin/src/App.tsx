import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth, AuthProvider, ProtectedRoute } from '@smexpress/auth';
import Layout from './components/Layout';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import UsersPage from './pages/UsersPage';
import FranchisesPage from './pages/FranchisesPage';
import CountriesPage from './pages/CountriesPage';
import FeatureFlagsPage from './pages/FeatureFlagsPage';
import RolesPage from './pages/RolesPage';
import SystemConfigPage from './pages/SystemConfigPage';
import CustomersPage from './pages/CustomersPage';
import CustomerDetailPage from './pages/CustomerDetailPage';
import AddressLookupPage from './pages/AddressLookupPage';

function LoginGuard() {
  const { isAuthenticated } = useAuth();
  if (isAuthenticated) return <Navigate to="/" replace />;
  return <LoginPage />;
}

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<LoginGuard />} />
        <Route element={<ProtectedRoute fallback={<Navigate to="/login" replace />}><Layout /></ProtectedRoute>}>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/franchises" element={<FranchisesPage />} />
          <Route path="/countries" element={<CountriesPage />} />
          <Route path="/feature-flags" element={<FeatureFlagsPage />} />
          <Route path="/roles" element={<RolesPage />} />
          <Route path="/system-config" element={<SystemConfigPage />} />
          <Route path="/customers" element={<CustomersPage />} />
          <Route path="/customers/:id" element={<CustomerDetailPage />} />
          <Route path="/address-lookup" element={<AddressLookupPage />} />
        </Route>
      </Routes>
    </AuthProvider>
  );
}
