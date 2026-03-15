import { ErrorBoundary } from '@/components/ErrorBoundary';
import { TokenInputModal } from '@/components/auth/TokenInputModal';
import { Header } from '@/components/layout/Header';
import { QueryProvider } from '@/components/providers/QueryProvider';
import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Rx',
  description: 'Plan-first training log system',
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
