export default async function Page({ params }: { params: { slug: string } }) {
  const accountId = params.slug;
  async function fetchTransactions(email: string, accountId: string) {}
  return <div>Account: {params.slug}</div>;
}
