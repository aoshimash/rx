import { Button } from '@/components/ui/button';
import { Dumbbell, Home, Settings } from 'lucide-react';
import Link from 'next/link';

/**
 * Navigation header with links to main pages
 */
export function Header() {
  return (
    <header className="border-b">
      <div className="container mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-6">
            <Link href="/" className="text-xl font-bold">
              Rx
            </Link>
            <nav className="flex items-center gap-2">
              <Link href="/">
                <Button variant="ghost" size="sm">
                  <Home className="h-4 w-4 mr-2" />
                  Week View
                </Button>
              </Link>
              <Link href="/plans">
                <Button variant="ghost" size="sm">
                  <Dumbbell className="h-4 w-4 mr-2" />
                  Plans
                </Button>
              </Link>
            </nav>
          </div>
          <Link href="/settings">
            <Button variant="ghost" size="sm">
              <Settings className="h-4 w-4 mr-2" />
              Settings
            </Button>
          </Link>
        </div>
      </div>
    </header>
  );
}
