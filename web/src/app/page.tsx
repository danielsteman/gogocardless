import Head from 'next/head';

export default function Home() {
  return (
    <div className='min-h-screen flex flex-col bg-gray-900 text-gray-200'>
      <Head>
        <title>PayPredict</title>
        <meta
          name='description'
          content='Predict your payments with PayPredict'
        />
        <link rel='icon' href='/favicon.ico' />
      </Head>

      {/* Header */}
      <header className='bg-gray-800 text-white p-4 shadow-lg'>
        <div className='container mx-auto flex justify-between items-center'>
          <h1 className='text-2xl font-bold'>PayPredict</h1>
          <nav>
            <a href='#features' className='px-4 hover:text-gray-400'>
              Features
            </a>
            <a href='#about' className='px-4 hover:text-gray-400'>
              About
            </a>
            <a href='#contact' className='px-4 hover:text-gray-400'>
              Contact
            </a>
            <a href='login' className='px-4 hover:text-gray-400'>
              Login
            </a>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <section className='bg-gray-800 text-white py-20'>
        <div className='container mx-auto text-center'>
          <h2 className='text-4xl font-bold mb-4'>Welcome to PayPredict</h2>
          <p className='text-xl mb-8'>
            Your smart solution to predict and manage payments efficiently.
          </p>
          <a
            href='#features'
            className='bg-blue-600 text-white px-6 py-3 rounded-full shadow-lg hover:bg-blue-500'
          >
            Learn More
          </a>
        </div>
      </section>

      {/* Features Section */}
      <section id='features' className='py-20'>
        <div className='container mx-auto text-center'>
          <h2 className='text-3xl font-bold mb-12'>Features</h2>
          <div className='flex flex-wrap justify-center'>
            <div className='w-full md:w-1/3 p-4'>
              <div className='bg-gray-800 p-6 rounded-lg shadow-lg'>
                <h3 className='text-2xl font-bold mb-4'>Payment Predictions</h3>
                <p>
                  Accurately predict upcoming payments and manage your finances
                  with ease.
                </p>
              </div>
            </div>
            <div className='w-full md:w-1/3 p-4'>
              <div className='bg-gray-800 p-6 rounded-lg shadow-lg'>
                <h3 className='text-2xl font-bold mb-4'>Detailed Insights</h3>
                <p>
                  Gain insights into your spending patterns and optimize your
                  budget.
                </p>
              </div>
            </div>
            <div className='w-full md:w-1/3 p-4'>
              <div className='bg-gray-800 p-6 rounded-lg shadow-lg'>
                <h3 className='text-2xl font-bold mb-4'>PSD2 Compliant</h3>
                <p>
                  Secure and compliant with the latest PSD2 regulations to
                  ensure your data&apos;s safety.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* About Section */}
      <section id='about' className='bg-gray-800 py-20'>
        <div className='container mx-auto text-center'>
          <h2 className='text-3xl font-bold mb-12'>About Us</h2>
          <p className='text-xl'>
            At PayPredict, we are dedicated to providing the best financial
            prediction tools to help you manage your payments and budget
            effectively. Our advanced algorithms and user-friendly interface
            make financial planning simple and accurate.
          </p>
        </div>
      </section>

      {/* Contact Section */}
      <section id='contact' className='py-20'>
        <div className='container mx-auto text-center'>
          <h2 className='text-3xl font-bold mb-12'>Contact Us</h2>
          <p className='text-xl mb-8'>
            Have any questions? Feel free to reach out to us.
          </p>
          <a
            href='mailto:contact@paypredict.com'
            className='bg-blue-600 text-white px-6 py-3 rounded-full shadow-lg hover:bg-blue-500'
          >
            Email Us
          </a>
        </div>
      </section>

      {/* Footer */}
      <footer className='bg-gray-800 text-white py-4'>
        <div className='container mx-auto text-center'>
          &copy; {new Date().getFullYear()} PayPredict. All rights reserved.
        </div>
      </footer>
    </div>
  );
}
