interface IComparingStats {
  name: string;
  total: number;
  failure: number;
  success: number;
}

export function ComparingStats({ total, failure, success, name }: IComparingStats) {

  const successPercentage = total === 0 ? '0' : ((success / total) * 100).toFixed(2);
  const failurePercentage = total === 0 ? '0' : ((failure / total) * 100).toFixed(2);

  return (
    <div>
      <h1 className="font-semibold text-gray-800 text-2xl">{name}</h1>
      <dl className="mt-5 grid grid-cols-1 divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow md:grid-cols-3 md:divide-x md:divide-y-0">
        <div className="px-4 py-5 sm:p-6">
          <dt className="text-base font-normal text-gray-900">Total</dt>
          <dd className="mt-1 flex items-baseline justify-between md:block lg:flex">
            <div className="flex items-baseline text-2xl font-semibold text-gray-800">
              {total}
            </div>
          </dd>
        </div>

        <div className="px-4 py-5 sm:p-6">
          <dt className="text-base font-normal text-gray-900">Success</dt>
          <dd className="mt-1 flex items-baseline justify-between md:block lg:flex">
            <div className="flex items-baseline text-2xl font-semibold text-gray-800">
              {success}
            </div>
            <div className="bg-green-100 text-green-800 inline-flex items-baseline rounded-full px-2.5 py-0.5 text-sm font-medium md:mt-2 lg:mt-0">
              {successPercentage} %
            </div>
          </dd>
        </div>

        <div className="px-4 py-5 sm:p-6">
          <dt className="text-base font-normal text-gray-900">Failure</dt>
          <dd className="mt-1 flex items-baseline justify-between md:block lg:flex">
            <div className="flex items-baseline text-2xl font-semibold text-gray-800">
              {failure}
            </div>
            <div className="bg-red-100 text-red-800 inline-flex items-baseline rounded-full px-2.5 py-0.5 text-sm font-medium md:mt-2 lg:mt-0">
              {failurePercentage} %
            </div>
          </dd>
        </div>

      </dl>
    </div>
  );
}
