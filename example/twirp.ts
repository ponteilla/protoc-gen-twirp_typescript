
interface TwirpErrorJSON {
  code: string;
  msg: string;
  meta: {[index:string]: string};
}

class TwirpError extends Error {
  code: string;
  meta: {[index:string]: string};

  constructor(te: TwirpErrorJSON) {
    super(te.msg);

    this.code = te.code;
    this.meta = te.meta;
  }
}

export const throwTwirpError = (resp: Response) => {
  return resp.json().then((err: TwirpErrorJSON) => { throw new TwirpError(err); })
};

export interface RequestParameters {
  baseUrl?: string;
  options?: any;
  fetch?: Fetch;
}

export const createRequest = (url: string, body: object, options?: any): Request => {
  const defaultOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(body),
  };

  const newOptions = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...(options && options.headers)
    },
  };

  return new Request(url, newOptions);
};

export type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>;
