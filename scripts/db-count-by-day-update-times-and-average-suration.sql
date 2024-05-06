-- count by day update and avg duration with specified workflow_id
    SELECT
        to_char(date, 'YYYY-MM-DD') AS Period,
        COUNT(workflow_logs.workflow_id) AS Count,
        AVG(duration) AS AvgDuration
    FROM
        generate_series(
        date_trunc('day', '2024-04-01'::date),
        date_trunc('day', now()),
        interval '1 day'
    ) AS date
    LEFT JOIN
        workflow_logs ON date_trunc('day', workflow_logs.created_at)::date = date::date
                     AND workflow_logs.state = 'completed'
                     AND workflow_logs.workflow_id = '94834729' -- 93324253 or 94834729
    GROUP BY
        date
    ORDER BY
        date;


-- count by day update and avg duration all workflow in project 
    SELECT
        to_char(date, 'YYYY-MM-DD') AS Period,
        COUNT(workflow_logs.workflow_id) AS Count,
        AVG(duration) AS AvgDuration
    FROM
        generate_series(
        date_trunc('day', now() - interval '34 day'),
        date_trunc('day', now()),
        interval '1 day'
    ) AS date
    LEFT JOIN
        workflow_logs ON date_trunc('day', workflow_logs.created_at)::date = date::date
                     AND workflow_logs.state = 'completed'
    GROUP BY
        date
    ORDER BY
        date;
