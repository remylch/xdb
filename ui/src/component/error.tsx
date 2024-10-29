import React from 'react';

type ErrorType = 'external' | 'internal';

type BaseError = {
  text: string;
  errorType: ErrorType;
};

type ExternalError = BaseError & {
  errorType: 'external';
  retry: () => void;
};

type InternalError = BaseError & {
  errorType: 'internal';
};

export const Error: React.FC<ExternalError | InternalError> = ({ errorType, text, ...props }) => {
  return (
    <span className="inline-flex items-center gap-x-1 rounded-md bg-red-100 px-2 py-1 text-xs font-medium text-red-600 w-fit">
      {text}
      {errorType === 'external' && 'retry' in props ? (
        <button onClick={props.retry} type="button" className="group relative -mr-1 h-3.5 w-3.5 rounded-sm hover:bg-red-500/20">
          <span className="sr-only">Remove</span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            className="h-4 w-4"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 0 0-3.7-3.7 48.678 48.678 0 0 0-7.324 0 4.006 4.006 0 0 0-3.7 3.7c-.017.22-.032.441-.046.662M19.5 12l3-3m-3 3-3-3m-12 3c0 1.232.046 2.453.138 3.662a4.006 4.006 0 0 0 3.7 3.7 48.656 48.656 0 0 0 7.324 0 4.006 4.006 0 0 0 3.7-3.7c.017-.22.032-.441.046-.662M4.5 12l3 3m-3-3-3 3"
            />
          </svg>
          <span className="absolute -inset-1" />
        </button>
      ) : null}
    </span>
  );
};
