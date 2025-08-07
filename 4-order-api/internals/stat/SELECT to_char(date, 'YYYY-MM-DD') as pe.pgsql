SELECT to_char(date, 'YYYY-MM-DD') as period, sum(clicks) FROM stats
WHERE date BETWEEN '01/01/2025' and '01/01/2026'
GROUP BY period
ORDER BY period

