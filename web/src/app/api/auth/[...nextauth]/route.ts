import NextAuth from "next-auth";
import GoogleProvider from "next-auth/providers/google";

export const options = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID as string,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET as string
    })
  ],
};

export const handler = NextAuth(options);

export { handler as GET, handler as POST };
