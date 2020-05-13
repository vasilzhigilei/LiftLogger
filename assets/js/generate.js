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