SELECT
  u.id,
  p.product_name,
  u.name
FROM
  users u
  JOIN products p on u.id = p.user_id
WHERE
  u.age > 25 and p.category = 'electronics'
GROUP BY
  u.id
ORDER BY
  p.price desc
LIMIT
  10