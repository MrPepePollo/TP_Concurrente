mtype = {LOCKED, UNLOCKED}

int resultado = 0; // Variable compartida para resultado
mtype mutex = UNLOCKED; // Mutex para controlar acceso a la variable compartida

proctype ParteRegresion() {
    int sumX = 0;
    int sumY = 0;
    int sumXY = 0;
    int sumX2 = 0;


    atomic {
        mutex == UNLOCKED -> 
        {
            mutex = LOCKED; // Bloquea el mutex para exclusión mutua
            sumX = sumX + 1;
            sumY = sumY + 1;
            sumXY = sumXY + 1;
            sumX2 = sumX2 + 1;

            mutex = UNLOCKED; // Libera el mutex para permitir otros procesos
        }
    }

    assert(mutex == UNLOCKED); // Confirma que el mutex está desbloqueado
}

// Proceso principal para correr varios procesos concurrentes
proctype Main() {
    run ParteRegresion();
    run ParteRegresion();
}

// Lanzar el proceso principal
init {
    run Main();
}
