import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Enable experimental features for Next.js 16
  experimental: {
    // Turbopack configuration
    turbo: {
      resolveAlias: {
        '@': '.',
        '@components': './components',
        '@lib': './lib',
        '@hooks': './hooks',
        '@stores': './stores',
        '@types': './types',
      },
    },
  },
  
  // Path aliases (for tsconfig.json paths to work)
  webpack: (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      '@': '.',
      '@components': './components',
      '@lib': './lib',
      '@hooks': './hooks',
      '@stores': './stores',
      '@types': './types',
    };
    return config;
  },
  
  // Image optimization
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**.githubusercontent.com',
      },
      {
        protocol: 'https',
        hostname: 'avatars.githubusercontent.com',
      },
    ],
  },
  
  // Headers for security
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
        ],
      },
    ];
  },
};

export default nextConfig;
