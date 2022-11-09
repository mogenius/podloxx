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
  statsTotalService?: {
    method?: RequestMethodEnum;
    endPoint?: string;
    header?: any;
  };
  statsFlowService?: {
    method?: RequestMethodEnum;
    endPoint?: string;
    header?: any;
  };
}

export const baseEnvironment: IEnvironment = {
  stage: 'dev',
  production: false,
  version: window.appVersion ?? pkg.version,
  baseUrl: 'http://127.0.0.1:1337',
  statsTotalService: {
    method: RequestMethodEnum.GET,
    endPoint: '/traffic/total',
    header: {
      // authorization: OVER AUTH INTERCEPTOR,
      contentType: 'application/json'
    }
  },
  statsFlowService: {
    method: RequestMethodEnum.GET,
    endPoint: '/traffic/flow',
    header: {
      // authorization: OVER AUTH INTERCEPTOR,
      contentType: 'application/json'
    }
  }
};
