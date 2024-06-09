import { getServerSession } from 'next-auth';
import { redirect } from 'next/navigation';
import authOptions from '../auth';
import SignIn from './SignIn';

const SignInPage = async () => {
  const session = await getServerSession(authOptions);

  if (session) {
    redirect('/banks');
  }

  return <SignIn />;
};

export default SignInPage;
