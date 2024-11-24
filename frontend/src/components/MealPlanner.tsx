import React, { useState, useEffect } from 'react';
import { Tag, X } from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';

interface MealItem {
  name: string;
  quantity: number;
  unit: string;
}

interface Meal {
  items: MealItem[];
}

interface MealPlan {
  [key: string]: Meal;
}

const MealPlanner: React.FC = () => {
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [availableIngredients, setAvailableIngredients] = useState<string[]>([]);
  const [currentInput, setCurrentInput] = useState('');
  const [calories, setCalories] = useState('2000');
  const [mealPlan, setMealPlan] = useState<MealPlan | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchIngredients();
  }, []);

  const fetchIngredients = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/ingredients');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setAvailableIngredients(data || []); // Assicuriamoci che sia sempre un array
    } catch (err) {
      console.error('Error fetching ingredients:', err);
      setError('Errore nel caricamento degli ingredienti');
      setAvailableIngredients([]); // Fallback a array vuoto in caso di errore
    }
  };

  const handleAddIngredient = () => {
    if (currentInput.trim() && !ingredients.includes(currentInput.trim().toLowerCase())) {
      setIngredients([...ingredients, currentInput.trim().toLowerCase()]);
      setCurrentInput('');
    }
  };

  const removeIngredient = (index: number) => {
    setIngredients(ingredients.filter((_, i) => i !== index));
  };

  const generateMealPlan = async () => {
    if (ingredients.length === 0) {
      setError('Seleziona almeno un ingrediente');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch('http://localhost:8080/api/generate-plan', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify({
          ingredients,
          targetCalories: parseInt(calories)
        })
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setMealPlan(data);
    } catch (err) {
      console.error('Error generating meal plan:', err);
      setError('Errore nella generazione del piano alimentare');
      setMealPlan(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <div className="bg-white rounded-lg shadow-lg p-6">
        <h1 className="text-2xl font-bold mb-6">Piano Alimentare Personalizzato</h1>

        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Ingredienti disponibili
          </label>
          <div className="flex gap-2">
            <Input
              list="ingredients"
              value={currentInput}
              onChange={(e) => setCurrentInput(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleAddIngredient()}
              placeholder="Inserisci un ingrediente..."
            />
            <Button onClick={handleAddIngredient}>
              Aggiungi
            </Button>
          </div>
          {availableIngredients.length > 0 && (
            <datalist id="ingredients">
              {availableIngredients.map((ing) => (
                <option key={ing} value={ing} />
              ))}
            </datalist>
          )}
        </div>

        {ingredients.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-6">
            {ingredients.map((ing, index) => (
              <span
                key={index}
                className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
              >
                <Tag className="w-4 h-4 mr-1" />
                {ing}
                <button
                  onClick={() => removeIngredient(index)}
                  className="ml-2 text-blue-600 hover:text-blue-800"
                >
                  <X className="w-4 h-4" />
                </button>
              </span>
            ))}
          </div>
        )}

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Calorie giornaliere
          </label>
          <Input
            type="number"
            value={calories}
            onChange={(e) => setCalories(e.target.value)}
            className="w-40"
            min="1000"
            max="5000"
          />
        </div>

        <Button
          onClick={generateMealPlan}
          disabled={ingredients.length === 0 || loading}
          className="w-full"
        >
          {loading ? (
            <span className="flex items-center justify-center">
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Generazione in corso...
            </span>
          ) : (
            'Genera Piano Alimentare'
          )}
        </Button>

        {mealPlan && Object.keys(mealPlan).length > 0 && (
          <div className="mt-6 space-y-4">
            {Object.entries(mealPlan).map(([mealType, meal]) => (
              <div key={mealType} className="bg-gray-50 rounded-lg p-4">
                <h3 className="text-lg font-medium capitalize mb-2">
                  {mealType}
                </h3>
                <div className="space-y-2">
                  {meal.items && meal.items.map((item: MealItem, idx: number) => (
                    <div key={idx} className="flex justify-between text-sm">
                      <span>{item.name}</span>
                      <span>{item.quantity} {item.unit}</span>
                    </div>
                  ))}
                  {(!meal.items || meal.items.length === 0) && (
                    <div className="text-sm text-gray-500">
                      Nessun alimento selezionato per questo pasto
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default MealPlanner;
