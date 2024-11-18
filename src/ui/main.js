import { searchStreet, saveClient } from './functions.js';

// Adiciona os eventos aos botões
document.getElementById('searchStreetBtn').addEventListener('click', searchStreet);
document.getElementById('saveClientBtn').addEventListener('click', saveClient);
