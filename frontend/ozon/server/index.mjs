'use strict';

import express from 'express';
import path from 'path';
import {fileURLToPath} from 'url';

const app = express();
app.use(express.json());
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const __dirnameOfProject = path.dirname(__dirname);

app.use(express.static(`${__dirnameOfProject}/src`))
app.get(/.*$/, (request, response) => {
    response.sendFile(`${__dirnameOfProject}/src/index.html`);
});

app.listen(3000);
