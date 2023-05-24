import { Client } from 'pg';

export default async function handler(request, response) {
  const client = new Client({
    connectionString: process.env.POSTGRES_URL,
    ssl: {
      rejectUnauthorized: false
    }
  });

  await client.connect();

  const createTableQuery = `
    CREATE TABLE IF NOT EXISTS screentime (
      id SERIAL PRIMARY KEY,
      screenshot TEXT,
      timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
  `;

  try {
    await client.query(createTableQuery);

    // Insert test data
    const testScreenshot = "This is a test screenshot.";
    const insertDataQuery = `
      INSERT INTO screentime (screenshot)
      VALUES ($1)
      RETURNING id;
    `;
    const result = await client.query(insertDataQuery, [testScreenshot]);

    // Query test data
    const selectDataQuery = `
      SELECT *
      FROM screentime
      WHERE id = $1;
    `;
    const data = await client.query(selectDataQuery, [result.rows[0].id]);

    response.status(200).json({
      body: data.rows[0],
      query: request.query,
      cookies: request.cookies,
    });
  } catch (error) {
    response.status(500).json({ error: 'Error processing request' });
  } finally {
    await client.end();
  }
}
