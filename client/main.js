import unfetch from "isomorphic-unfetch";

import '@fortawesome/fontawesome-free-webfonts';
import '@fortawesome/fontawesome-free-webfonts/css/fa-solid.css';

const Entry = ({ kind, title, authors }) => `
<li class="entry ${kind}">
  <i class="fas fa-${kind=='book'?'book':'compact-disc'}"></i>
  <strong class="title">${title}</strong>
  <span class="author">${(authors||[]).join(', ')}</span>
</li>
`;

const API_URL = process.env["API_URL"] || 'http://localhost:8080';

document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('search');
  const result = document.getElementById('result');
  const took = document.getElementById('took');

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const button = e.target.submit;
    button.disabled = true;

    try {
      const q = e.target.q.value;
      const resp = await unfetch(`${API_URL}?q=${q}`)
      const body = await resp.json();

      const entries = body.data.map(Entry);
      result.innerHTML = `<ul>` + entries.join('') + `</ul>`
      took.innerHTML = `Took ${body.took / 1000000000} seconds`;

    } finally {
      button.disabled = false;
    }
  });

  form.q.focus();
});
