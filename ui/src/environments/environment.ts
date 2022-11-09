import { baseEnvironment } from '@lox/environments/base.environment';
import * as deepmerge from 'deepmerge';

// const url = 'https://api.dev.mogenius.com'; // TODO Austauschen
const url = 'http://127.0.0.1:1337';

export const environment = deepmerge(baseEnvironment, {
  stage: 'local',
  production: false,
  baseUrl: url,
  statsService: {}
});
