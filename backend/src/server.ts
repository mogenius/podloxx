process.env['NODE_CONFIG_DIR'] = __dirname + '/configs';

import 'dotenv/config';
import App from '@/app';
import CadvisorRoute from '@routes/cadvisor.route';
import validateEnv from '@utils/validateEnv';

validateEnv();

const app = new App([new CadvisorRoute()]);

app.listen();
