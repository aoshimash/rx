import type { Metadata } from 'next';
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
      <body>{children}</body>
    </html>
  );
}
