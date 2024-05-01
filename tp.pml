// Definición de la estructura de datos para mantener resultados parciales
mtype = {LOCKED, UNLOCKED}

int resultado = 0; // Variable compartida para resultado
mtype mutex = UNLOCKED; // Mutex para controlar acceso a la variable compartida

// Proceso para simular parte de la regresión que modifica una variable compartida
proctype ParteRegresion() {
    int sumX = 0;
    int sumY = 0;
    int sumXY = 0;
    int sumX2 = 0;

    // Simulación de trabajo en concurrencia
    atomic {
        // Asegura que solo un proceso puede entrar al bloque crítico
        mutex == UNLOCKED -> 
        {
            mutex = LOCKED; // Bloquea el mutex para exclusión mutua

            // Sección crítica: modificaciones a la variable compartida
            sumX = sumX + 1;
            sumY = sumY + 1;
            sumXY = sumXY + 1;
            sumX2 = sumX2 + 1;

            mutex = UNLOCKED; // Libera el mutex para permitir otros procesos
        }
    }

    // Verificar que el mutex se liberó correctamente
    assert(mutex == UNLOCKED); // Confirma que el mutex está desbloqueado
}

// Proceso principal para correr varios procesos concurrentes
proctype Main() {
    // Lanzar dos procesos para simular concurrencia
    run ParteRegresion();
    run ParteRegresion();
}

// Lanzar el proceso principal
init {
    run Main();
}

