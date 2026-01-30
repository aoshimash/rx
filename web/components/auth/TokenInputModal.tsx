'use client';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { authStore } from '@/stores/auth';
import { useEffect, useState } from 'react';

/**
 * Token input modal for lazy authentication
 *
 * Triggered when API operations require authentication
 */
export function TokenInputModal() {
  const [open, setOpen] = useState(false);
  const [token, setToken] = useState('');

  useEffect(() => {
    const handleAuthRequired = () => {
      setOpen(true);
    };

    window.addEventListener('auth:required', handleAuthRequired);
    return () => window.removeEventListener('auth:required', handleAuthRequired);
  }, []);

  const handleSave = () => {
    if (token.trim()) {
      authStore.setToken(token.trim());
      setOpen(false);
      setToken('');
      // Reload the page to retry failed requests
      window.location.reload();
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Authentication Required</DialogTitle>
          <DialogDescription>
            Please enter your API authentication token to continue.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="auth-token">API Token</Label>
            <Input
              id="auth-token"
              type="password"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder="Enter your token"
              onKeyDown={(e) => e.key === 'Enter' && handleSave()}
            />
            <p className="text-sm text-muted-foreground">You can also set this in Settings page.</p>
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={!token.trim()}>
              Save Token
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
