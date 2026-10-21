-- exercises が body_parts を参照しているので、種目から先に消す
DELETE FROM exercises WHERE id IN (
    'bench-press', 'incline-dumbbell-press', 'cable-fly',
    'deadlift', 'lat-pulldown', 'bent-over-row',
    'squat', 'leg-press', 'leg-curl',
    'shoulder-press', 'side-raise',
    'dumbbell-curl', 'triceps-pushdown',
    'crunch', 'leg-raise'
);

DELETE FROM body_parts WHERE code IN ('chest', 'back', 'legs', 'shoulders', 'arms', 'abs');
