import { RequestMethodEnum } from '@lox/enums/request-method.enum';
import pkg from '../../package.json';

declare global {
  interface Window {
    appVersion: string;
  }
}

export interface IEnvironment {
  stage: 'dev' | 'prod' | 'production' | 'local';
  production: boolean;
  version?: string;
  //! SERVICES ...
  baseUrl: string;
  statsService?: {
    method?: RequestMethodEnum;
    endPoint?: string;
    header?: any;
  };
  cAdvisorService: {
    baseUrl: string;
    cAdvisorRawData: {
      method?: RequestMethodEnum;
      endPoint?: string;
      header?: any;
    };
  };
}

export const baseEnvironment: IEnvironment = {
  stage: 'dev',
  production: false,
  version: window.appVersion ?? pkg.version,
  baseUrl: 'https://api.dev.mogenius.com', // TODO Change
  statsService: {
    method: RequestMethodEnum.GET,
    endPoint: '/service/network-stats',
    header: {
      // authorization: OVER AUTH INTERCEPTOR,
      contentType: 'application/json'
    }
  },

  cAdvisorService: {
    baseUrl: 'http://127.0.0.1:8001', // TODO Change
    cAdvisorRawData: {
      method: RequestMethodEnum.GET,
      endPoint: '/api/v1/nodes/aks-devpool2-28911264-vmss000000/proxy/metrics/cadvisor',
      header: {
        // authorization: OVER AUTH INTERCEPTOR,
        contentType: 'text/plain'
      }
    }
  }
};
