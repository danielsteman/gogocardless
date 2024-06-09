import { getToken } from 'next-auth/jwt';
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const secret = process.env.NEXTAUTH_SECRET;

export async function middleware(req: NextRequest) {
  const token = await getToken({ req, secret });

  // Define the protected routes
  const protectedRoutes = ['/banks', '/account', '/settings'];

  // Check if the current route is protected
  const isProtectedRoute = protectedRoutes.some((route) => req.nextUrl.pathname.startsWith(route));

  if (isProtectedRoute) {
    if (!token) {
      // Redirect to login if user is not authenticated
      const url = req.nextUrl.clone();
      url.pathname = '/login';
      return NextResponse.redirect(url);
    }
  }

  return NextResponse.next();
}

// Configures the middleware to run on specific paths
export const config = {
  matcher: ['/banks/:path*', '/account/:path*', '/settings/:path*'],
};
