import type { Metadata } from 'next';
import { QueryProvider } from '@/components/providers/QueryProvider';
import { TokenInputModal } from '@/components/auth/TokenInputModal';
import { Header } from '@/components/layout/Header';
import { ErrorBoundary } from '@/components/ErrorBoundary';
import './globals.css';

export const metadata: Metadata = {
  title: 'OPTel Workout',
  description: 'Workout management for the human physical layer',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <ErrorBoundary>
          <QueryProvider>
            <Header />
            {children}
            <TokenInputModal />
          </QueryProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
