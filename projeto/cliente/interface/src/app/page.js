'use client';

import { useEffect, useState } from 'react';
import { GetRotas } from '@/app/api/rotas'; // Importa corretamente

export default function Home() {
  const [rotas, setRotas] = useState(null);

  useEffect(() => {
    // Chama a função GetRotas ao montar o componente
    GetRotas()
      .then(data => {
        if (data) {
          setRotas(data);
          console.log(data); // Exibe os dados no console
        }
      })
      .catch(error => {
        console.error('Erro ao buscar rotas:', error);
      });
  }, []);

  return (
    <>
      {rotas ? (
        <pre>{JSON.stringify(rotas, null, 2)}</pre> // Mostra as rotas em JSON formatado
      ) : (
        <p>Carregando rotas...</p>
      )}
    </>
  );
}
