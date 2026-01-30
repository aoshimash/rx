'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { authStore } from '@/stores/auth';
import { useEffect, useState } from 'react';

export default function SettingsPage() {
  const [token, setToken] = useState('');
  const [hasToken, setHasToken] = useState(false);

  useEffect(() => {
    setHasToken(authStore.hasToken());
  }, []);

  const handleSaveToken = () => {
    if (token.trim()) {
      authStore.setToken(token.trim());
      setHasToken(true);
      setToken('');
    }
  };

  const handleClearToken = () => {
    authStore.clearToken();
    setHasToken(false);
    setToken('');
  };

  return (
    <main className="container mx-auto p-6 max-w-2xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-1">Configure your application preferences</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>API Authentication</CardTitle>
          <CardDescription>
            Manage your API authentication token for backend access.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="api-token">API Token</Label>
            <Input
              id="api-token"
              type="password"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder={hasToken ? '••••••••' : 'Enter your API token'}
            />
            <p className="text-sm text-muted-foreground">
              This token is stored locally in your browser and used to authenticate API requests.
            </p>
          </div>

          <div className="flex gap-2">
            <Button onClick={handleSaveToken} disabled={!token.trim()}>
              Save Token
            </Button>
            {hasToken && (
              <Button variant="destructive" onClick={handleClearToken}>
                Clear Token
              </Button>
            )}
          </div>

          {hasToken && <div className="text-sm text-green-600">✓ Token is configured</div>}
        </CardContent>
      </Card>
    </main>
  );
}
