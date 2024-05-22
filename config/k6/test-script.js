import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.0.0/index.js';

export let options = {
    vus: 1000,
    duration: '60s',
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% das requisições devem ser completadas em menos de 500ms
    },
};

const BASE_URL = 'http://app:8080';
let userIds = []; // Array para armazenar os userIds

export default function () {
    // Step 1: Create wallets and deposit funds (executed less frequently)
    if (Math.random() < 0.1) { // 10% chance to execute this block
        group('Create wallets and deposit', function () {
            for (let i = 0; i < 2; i++) {
                let walletRes = http.post(`${BASE_URL}/wallets`);
                check(walletRes, {
                    'wallet creation is status 200': (r) => r.status === 200,
                });

                let walletData = JSON.parse(walletRes.body);
                let userId = walletData.user_id;
                userIds.push(userId); // Armazena o userId criado

                let depositRes = http.patch(`${BASE_URL}/wallets/deposit`, JSON.stringify({
                    userId: userId,
                    amount: Math.random() * 1000, // valor aleatório para o depósito
                    vibranium: Math.floor(Math.random() * 100), // valor aleatório de vibranium
                }), {
                    headers: { 'Content-Type': 'application/json' },
                });
                check(depositRes, {
                    'deposit is status 200': (r) => r.status === 200,
                });
            }
        });
    }

    // Step 2: Create buy/sell orders (executed more frequently)
    if (userIds.length > 0) { // Verifica se há userIds disponíveis
        group('Create buy/sell orders', function () {
            for (let i = 0; i < 10; i++) { // Increased number of /orders requests
                let userId = userIds[Math.floor(Math.random() * userIds.length)]; // Seleciona um userId aleatório do array
                let orderType = Math.random() > 0.5 ? 1 : 2; // 1 para compra, 2 para venda
                let amount = Math.floor(Math.random() * 100); // quantidade de vibranium
                let price = Math.random() * 100; // preço por unidade de vibranium

                let orderRes = http.post(`${BASE_URL}/orders`, JSON.stringify({
                    userId: userId,
                    type: orderType,
                    amount: amount,
                    price: price,
                }), {
                    headers: { 'Content-Type': 'application/json' },
                });
                check(orderRes, {
                    'order creation is status 200': (r) => r.status === 200,
                });
            }
        });
    }

    sleep(1); // Pequena pausa entre as requisições
}
