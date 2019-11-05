
import * as Service from './service.pb'
import * as Twirp from './twirp'

export const makeHat = (requestParams: Twirp.RequestParameters, size: Service.Size): Promise<Service.Hat> => {
  const url = requestParams.baseUrl + "/example.Haberdasher/MakeHat";
  const body = Service.SizeToJSON(size);
  const fetchRequest: Twirp.Fetch = requestParams.fetch ? requestParams.fetch : window.fetch.bind(window);
  
  return fetchRequest(Twirp.createRequest(url, body, requestParams.options)).then((resp) => {
    if(!resp.ok) {
      return Twirp.throwTwirpError(resp);
    }

    return resp.json().then(Service.JSONToHat);
  });
};
