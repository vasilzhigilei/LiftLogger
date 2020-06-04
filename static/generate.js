/**
 * Generates 1 rep max calculation for given pounds & reps
 * utilizes the Wathen formula
 * @param pounds - weight lifted for compound lift
 * @param reps - repetitions given compound lift was pressed or pulled
 * @returns {number} - the one rep max estimate
 */
function generate1RM(pounds, reps) {
    return (100 * pounds) / (48.8 + 53.8 * Math.E**(-0.075 * reps))
}

/**
 * Generates X rep max given a 1 rep max
 * @param max - weight lifted for 1 rep max
 * @param reps - reps to convert the aforementioned pounds to
 */
function generateXRM(max, reps) {
    return (max * (48.8 + 53.8 * Math.E**(-0.075 * reps))) / 100
}