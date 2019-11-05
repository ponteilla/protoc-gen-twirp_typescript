package json

const messageTemplate = `
{{- range .Enums }}
export enum {{ EntityNamer . }} {
  {{- range .Values }}
  {{ .Name }} = "{{ .Name }}",
  {{- end}}
}
{{ end }}
{{- range .Messages }}
{{- $message_name := EntityNamer . }}
export interface {{ $message_name }} {
  {{- range .Fields }}
  {{ FieldNamer . }}{{ if IsMessage . }}?{{ end }}: {{ FieldTyper . }}{{ if .Type.IsRepeated }}[]{{ end }};
  {{- end }}
}

export interface {{ $message_name }}JSON {
  {{- range .Fields }}
  {{ JSONFieldNamer . }}?: {{ JSONFieldTyper . }}{{ if .Type.IsRepeated }}[]{{ end }};
  {{- end }}
}

export const {{ $message_name }}ToJSON = (m: {{ $message_name }}): {{ $message_name }}JSON => {
  return {
    {{- range .Fields }}
    {{ JSONFieldNamer . }}: {{ JSONCaster . }},
    {{- end }}
  };
};

export const JSONTo{{ $message_name }} = (m: {{ $message_name }}JSON): {{ $message_name }} => {
  return {
    {{- range .Fields }}
    {{ FieldNamer . }}: {{ Caster .}},
    {{- end }}
  };
};
{{ end }}
{{- range .MapMessages }}
export interface {{ .Name }} {
  [key: string]: {{ .Type }};
}

export interface {{ .Name }}JSON {
  [key: string]: {{ .Type }}{{- if .IsMessage }}JSON{{ end }};
}

{{- if .IsMessage }}
export const JSONTo{{ .Name }} = (m: {{ .Name }}JSON): {{ .Name }} => {
  return Object.keys(m).reduce((acc, key) => {
    acc[key] = JSONTo{{ .Type }}(m[key]);
    return acc;
  }, {} as {{ .Name }});
};

export const {{ .Name }}ToJSON = (m: {{ .Name }}): {{ .Name }}JSON => {
  return Object.keys(m).reduce((acc, key) => {
    acc[key] = {{ .Type }}ToJSON(m[key]);
    return acc;
  }, {} as {{ .Name }}JSON);
};
{{ end }}
{{- end -}}
`

const clientTemplate = `
import * as Service from './service.pb'
import * as Twirp from './twirp'

{{- range .Services -}}
{{- $service_name := .Name }}
{{- range .Methods }}
{{ $input_name := EntityNamer .Input }}
export const {{ .Name.LowerCamelCase }} = (requestParams: Twirp.RequestParameters, {{ $input_name | LowerCamelCaser }}: Service.{{ $input_name}}): Promise<Service.{{ EntityNamer .Output }}> => {
  const url = requestParams.baseUrl + "/{{ .Package.ProtoName }}.{{ $service_name }}/{{ .Name }}";
  const body = Service.{{ $input_name }}ToJSON({{ $input_name | LowerCamelCaser }});
  const fetchRequest: Twirp.Fetch = requestParams.fetch ? requestParams.fetch : window.fetch.bind(window);

  return fetchRequest(Twirp.createRequest(url, body, requestParams.options)).then((resp) => {
    if(!resp.ok) {
      return Twirp.throwTwirpError(resp);
    }

    return resp.json().then(Service.JSONTo{{ EntityNamer .Output }});
  });
};
{{- end }}
{{- end }}
`

const twirpTemplate = `
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
`
