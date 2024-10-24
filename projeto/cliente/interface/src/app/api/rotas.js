export async function GetRotas() {
    try {
      const response = await fetch('http://localhost:8000/rota');
      
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      
      const data = await response.json(); // Converte a resposta em JSON
      console.log(data); // Manipula os dados recebidos
      
      return data; // Retorna os dados, se necessário
    } catch (error) {
      console.error('Houve um problema com a requisição:', error);
    }
  }
  