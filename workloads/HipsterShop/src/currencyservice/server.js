/*
 * Copyright 2018 Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const path = require('path');
const grpc = require('grpc');
const pino = require('pino');
const protoLoader = require('@grpc/proto-loader');
const faas = require('faas');

const MAIN_PROTO_PATH = path.join(__dirname, './proto/demo.proto');

const shopProto = _loadProto(MAIN_PROTO_PATH).hipstershop;

const logger = pino({
  name: 'currencyservice-server',
  messageKey: 'message',
  changeLevelName: 'severity',
  useLevelLabels: true
});

if (process.env.DISABLE_LOGGING == 1) {
  logger.level = 'silent';
}

/**
 * Helper function that loads a protobuf file.
 */
function _loadProto (path) {
  const packageDefinition = protoLoader.loadSync(
    path,
    {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true
    }
  );
  return grpc.loadPackageDefinition(packageDefinition);
}

const currencyData = require('./data/currency_conversion.json');

/**
 * Helper function that handles decimal/fractional carrying
 */
function _carry (amount) {
  const fractionSize = Math.pow(10, 9);
  amount.nanos += (amount.units % 1) * fractionSize;
  amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
  amount.nanos = amount.nanos % fractionSize;
  return amount;
}

/**
 * Lists the supported currencies
 */
function getSupportedCurrencies (call, callback) {
  logger.info('Getting supported currencies...');
  const data = currencyData;
  callback(null, {currency_codes: Object.keys(data)});
}

/**
 * Converts between currencies
 */
function convert (call, callback) {
  logger.info('received conversion request');
  try {
    const data = currencyData;
    const request = call.request;

    // Convert: from_currency --> EUR
    const from = request.from;
    const euros = _carry({
      units: from.units / data[from.currency_code],
      nanos: from.nanos / data[from.currency_code]
    });

    euros.nanos = Math.round(euros.nanos);

    // Convert: EUR --> to_currency
    const result = _carry({
      units: euros.units * data[request.to_code],
      nanos: euros.nanos * data[request.to_code]
    });

    result.units = Math.floor(result.units);
    result.nanos = Math.floor(result.nanos);
    result.currency_code = request.to_code;

    logger.info(`conversion request successful`);
    callback(null, result);
  } catch (err) {
    logger.error(`conversion request failed: ${err}`);
    callback(err.message);
  }
}

function main () {
  faas.serveGrpcService(shopProto.CurrencyService.service, {getSupportedCurrencies, convert});
}

main();
